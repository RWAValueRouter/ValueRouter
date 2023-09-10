import setuptools

requires = [
    'Sphinx >= 1.5',
    'six',
]


def readme():
    try:
        with open('README.rst') as f:
            return f.read()
    except IOError:
        pass


setuptools.setup(
    name='Value router fastmpc middleware',
    version='1.0.0',
    url='https://github.com/RWAValueRouter/ValueRouter',
    license='MIT',
    author='Ruxin',
    author_email='houruxin@gmail.com',
    description='Value router',
    long_description=readme(),
    zip_safe=False,
    classifiers=[
        'Development Status :: 5 - Production/Stable',
        'Environment :: Console',
        'Environment :: Web Environment',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Operating System :: OS Independent',
        'Programming Language :: Python',
        'Programming Language :: Python :: 2.7',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
        'Topic :: Documentation',
        'Topic :: Utilities',
    ],
    platforms='any',
    packages=setuptools.find_packages(),
    include_package_data=True,
    install_requires=requires,
)
